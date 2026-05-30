set -xe
TMP_PUB_KEY="$(openssl genrsa 4096 | tee id_rsa | openssl rsa -pubout | grep -v 'PUBLIC KEY' | tr -d '\n')"
CREATE_QUERY="CALL PROCEDURES.ORGADMIN.CREATE_PR_ACC($1,'$TMP_PUB_KEY');"
RESULT_JSON="$(snow --config-file config.toml sql --debug --user "$SNOWFLAKE_USER" --accountname "$SNOWFLAKE_ACCOUNT_NAME" --role "$SNOWFLAKE_ROLE" --private-key-file "$SNOWFLAKE_PRIVATE_KEY_PATH" -q "$CREATE_QUERY" --format json)"
echo $RESULT_JSON
PR_ACCOUNT_NAME=$(echo $RESULT_JSON | yq -r '.[0].ACCOUNT_NAME')
PR_USER_NAME=$(echo $RESULT_JSON | yq -r '.[0].ACCOUNT_NAME')
PR_ORG_NAME=$(echo $RESULT_JSON | yq -r '.[0].ORGANIZATION_NAME')
PR_HOST=$(echo $RESULT_JSON | yq -r '.[0].LOGIN_URL')
echo "Waiting for up to 5m for the account's url ($PR_HOST) to come live"
attempt_counter=0
max_attempts=30
until $(curl --output /dev/null --silent --fail https://$PR_HOST); do
    if [ ${attempt_counter} -eq ${max_attempts} ];then
      echo "timedout waiting for login url $PR_HOST to come live"
      exit 1
    fi

    printf '.'
    attempt_counter=$(($attempt_counter+1))
    sleep 10
done
snow --config-file config.toml sql -v --enable-diag --host "$PR_HOST" --user "$PR_USER_NAME" --accountname "${PR_ORG_NAME}-${PR_USER_NAME}" --private-key-file id_rsa --role ACCOUNTADMIN -q "ALTER USER SET WORKLOAD_IDENTITY = ( TYPE = OIDC ISSUER='https://token.actions.githubusercontent.com' SUBJECT='$PR_SUBJECT' );"

 yq --null-input ".PR_USER_NAME = \"$PR_USER_NAME\" | .PR_ACCOUNT_NAME = \"$PR_ACCOUNT_NAME\" | .PR_ORG_NAME = \"$PR_ORG_NAME\" | .PR_HOST = \"$PR_HOST\"" | tee pr-account.yaml