terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "2.15.0"
    }
  }
}

provider "snowflake" {
  role = "ORGADMIN" # Required to create an owner-privileges procedure that can create accounts
  alias = "orgadm"
  preview_features_enabled = toset([ 
    "snowflake_procedure_sql_resource"
   ])
}

provider "snowflake" {
  role = "ACCOUNTADMIN" # Required to create an owner-privileges procedure that can create accounts
  alias = "accadm"
}

variable "procedure_db" {
  type = string
  default = "PROCEDURES"
}

variable "procedure_schema" {
  type = string
  default = "ORGADMIN"
}

variable "user_pub_key_file" {
  type = string
  default = "./id_rsa.pub"
}

locals {
  pub_key_lines = split("\n",chomp(file(var.user_pub_key_file)))
  pub_key_extracted = join("",slice(local.pub_key_lines,1,length(local.pub_key_lines)-1))
}

resource "snowflake_database" "procedure_db" {
  provider = snowflake.accadm
  name = var.procedure_db
}

resource "snowflake_schema" "procedure_schema" {
  provider = snowflake.accadm
  name = var.procedure_schema
  database = snowflake_database.procedure_db.name
}

resource "snowflake_grant_privileges_to_account_role" "orgadm_create_procedure" {
  provider = snowflake.accadm
  account_role_name = "ORGADMIN"
  all_privileges = true
  on_schema_object {
    future {
      in_schema = snowflake_schema.procedure_schema.fully_qualified_name
      object_type_plural = "PROCEDURES"
    }
  }
}

resource "snowflake_procedure_sql" "create_acc" {
  provider = snowflake.orgadm
  name = "CREATE_PR_ACC"

  database = snowflake_database.procedure_db.name
  schema = snowflake_schema.procedure_schema.name

  execute_as = "OWNER"
  return_type = "TABLE()"
  arguments {
    arg_name = "PR_NUMBER"
    arg_data_type = "INT"
  }
  arguments {
    arg_name = "RSA_PUBLIC_KEY"
    arg_data_type = "VARCHAR(4096)"
  }
  procedure_definition = <<EOT
BEGIN
LET ACCOUNT_NAME VARCHAR(32) := CONCAT('PR_',TO_VARCHAR(:PR_NUMBER));
SHOW ACCOUNTS LIKE :ACCOUNT_NAME;
LET cnt INTEGER := (SELECT COUNT(*) FROM TABLE(RESULT_SCAN(LAST_QUERY_ID())) WHERE "account_name" = :ACCOUNT_NAME);
-- Only create account if it doesn't exist already
IF (cnt > 0) THEN
  -- Account exists, "reset" it, by renaming the account out of the way, dropping it and then creating a new account with that name
  LET NEW_NAME VARCHAR(128) := (SELECT RANDSTR(8,RANDOM()) || '_' || :ACCOUNT_NAME);

  DROP ACCOUNT IDENTIFIER(:ACCOUNT_NAME) GRACE_PERIOD_IN_DAYS=3;

  ALTER ACCOUNT IDENTIFIER(:ACCOUNT_NAME) RENAME TO IDENTIFIER( :NEW_NAME ) SAVE_OLD_URL=FALSE;
END IF;
  CREATE ACCOUNT IDENTIFIER(:ACCOUNT_NAME)
      EMAIL = 'provider-pr@bronze-deer.de'
      ADMIN_NAME = :ACCOUNT_NAME
      ADMIN_USER_TYPE = SERVICE
      REGION = 'AWS_EU_CENTRAL_1'

      ADMIN_RSA_PUBLIC_KEY = :RSA_PUBLIC_KEY
      EDITION = STANDARD
  ;
LET res RESULTSET := (SELECT :ACCOUNT_NAME AS ACCOUNT_NAME, CURRENT_ORGANIZATION_NAME() AS ORGANIZATION_NAME, LOWER(CURRENT_ORGANIZATION_NAME()) || '-' || LOWER(:ACCOUNT_NAME) || '.snowflakecomputing.com' AS LOGIN_URL);
RETURN TABLE(res);
END;

EOT
}

resource "snowflake_procedure_sql" "drop_acc" {
  provider = snowflake.orgadm
  name = "DROP_PR_ACC"

  database = snowflake_database.procedure_db.name
  schema = snowflake_schema.procedure_schema.name

  execute_as = "OWNER"
  return_type = "VARCHAR(4096)"
  arguments {
    arg_name = "PR_NUMBER"
    arg_data_type = "INT"
  }
  procedure_definition = <<EOT
BEGIN

LET ACCOUNT_NAME VARCHAR(32) := CONCAT('PR_',TO_VARCHAR(:PR_NUMBER));

SHOW ACCOUNTS LIKE :ACCOUNT_NAME;
LET ACCOUNT_EXISTS BOOLEAN := (SELECT COUNT(*) > 0 FROM TABLE(RESULT_SCAN(LAST_QUERY_ID())));

IF (NOT :ACCOUNT_EXISTS) THEN
  RETURN :ACCOUNT_NAME;
END IF;

LET NEW_NAME VARCHAR := (SELECT :ACCOUNT_NAME || '_OLD_' || RANDSTR(6, RANDOM()));
ALTER ACCOUNT IF EXISTS IDENTIFIER(:ACCOUNT_NAME) RENAME TO IDENTIFIER(:NEW_NAME) SAVE_OLD_URL=FALSE;
DROP ACCOUNT IF EXISTS IDENTIFIER(:NEW_NAME) GRACE_PERIOD_IN_DAYS = 3;
RETURN :ACCOUNT_NAME;
END;

EOT
}

resource "snowflake_account_role" "pr_acc_manager" {
  provider = snowflake.accadm
  name = "PR_ACC_MANAGER"
}

resource "snowflake_grant_privileges_to_account_role" "use_db" {
  provider = snowflake.accadm
  account_role_name = snowflake_account_role.pr_acc_manager.name
  privileges = [ "USAGE" ]
  on_account_object {
    object_name = snowflake_database.procedure_db.name
    object_type = "DATABASE"
  }
}

resource "snowflake_grant_privileges_to_account_role" "use_schema" {
  provider = snowflake.accadm
  account_role_name = snowflake_account_role.pr_acc_manager.name
  privileges = [ "USAGE" ]
  on_schema {
    schema_name = snowflake_schema.procedure_schema.fully_qualified_name
  }
}

resource "snowflake_grant_privileges_to_account_role" "grant_create_procedure" {
  account_role_name = snowflake_account_role.pr_acc_manager.name
  privileges = [ "USAGE" ]
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = snowflake_procedure_sql.create_acc.fully_qualified_name
  }
}

resource "snowflake_grant_privileges_to_account_role" "grant_drop_procedure" {
  account_role_name = snowflake_account_role.pr_acc_manager.name
  privileges = [ "USAGE" ]
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = snowflake_procedure_sql.drop_acc.fully_qualified_name
  }
}

resource "snowflake_service_user" "pr_acc_manager" {
  provider = snowflake.accadm
  name = "GH_ACTIONS_PR_ACC_MANAGER"
  default_role = snowflake_account_role.pr_acc_manager.name
  rsa_public_key = local.pub_key_extracted
  default_warehouse = snowflake_warehouse.pr_acc_manager_warehouse.name
}

resource "snowflake_grant_account_role" "grant_acc_mgr_role_to_user" {
  provider = snowflake.accadm
  role_name = snowflake_account_role.pr_acc_manager.name
  user_name = snowflake_service_user.pr_acc_manager.name
}

resource "snowflake_warehouse" "pr_acc_manager_warehouse" {
  provider = snowflake.accadm
  name = "GH_ACTIONS_PR_ACC_MANAGER_WH"
  warehouse_size = "X-SMALL"
}

resource "snowflake_grant_privileges_to_account_role" "grant_acc_mgr_wh_usage" {
  provider = snowflake.accadm
  account_role_name = snowflake_account_role.pr_acc_manager.name
  privileges = [ "USAGE" ]
  on_account_object {
    object_type = "WAREHOUSE"
    object_name = snowflake_warehouse.pr_acc_manager_warehouse.name
  }
}