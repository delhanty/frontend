APP_ENV = dev
CONSOLE_HANDLERS = auth address org scangroup user webdata

test:
	go test ./... -cover

build:
	$(foreach var,$(CONSOLE_HANDLERS),GOOS=linux go build -o dist/console/main ./cmd/console/$(var)/ && zip -j dist/console/$(var)_handler.zip dist/console/main && rm dist/console/main;)

upload:
	aws s3 sync dist/console/ s3://linkai-infra/frontend/lambdas/console/



# Authentication
buildauth:
	GOOS=linux go build -o dist/console/main ./cmd/console/auth/ && zip -j dist/console/auth_handler.zip dist/console/main && rm dist/console/main

deployauth: buildauth upload
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/console/auth_handler.zip --function-name dev-console-handler-auth

# Lambda Authorizer
buildauthorizer:
	GOOS=linux go build -o dist/authorizer/main ./cmd/authorizer/ && zip -j dist/authorizer/lambda_authorizer.zip dist/authorizer/main && rm dist/authorizer/main

uploadauthorizer:
	aws s3 sync dist/authorizer/ s3://linkai-infra/frontend/lambdas/authorizer/

deployauthorizer: buildauthorizer uploadauthorizer
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/authorizer/lambda_authorizer.zip --function-name dev-console-handler-lambda-authorizer

# Static Contents Authorizer
buildstaticauthorizer:
	GOOS=linux go build -o dist/staticauthorizer/main ./cmd/staticauthorizer/ && zip -j dist/staticauthorizer/static_authorizer.zip dist/staticauthorizer/main && rm dist/staticauthorizer/main

uploadstaticauthorizer: 
	aws s3 sync dist/staticauthorizer/ s3://linkai-infra/frontend/lambdas/staticauthorizer/

deploystaticauthorizer: buildstaticauthorizer uploadstaticauthorizer
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/staticauthorizer/static_authorizer.zip --function-name dev-console-handler-static-authorizer

# Admin features 
buildadmin:
	GOOS=linux go build -o dist/console/admin/main ./cmd/console/admin/ && zip -j dist/console/admin/admin_handler.zip dist/console/admin/main && rm dist/console/admin/main

uploadadmin: 
	aws s3 sync dist/console/admin s3://linkai-infra/frontend/lambdas/console/admin/

deployadmin: buildadmin uploadadmin
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/console/admin/admin_handler.zip --function-name dev-console-handler-admin

# Support Provision
supportprovision:
	docker build -t support_org_provision -f Dockerfile.support_org_provision .

pushsupportprovision: supportprovision
	docker tag support_org_provision:latest 447064213022.dkr.ecr.us-east-1.amazonaws.com/support_org_provision:latest && docker push 447064213022.dkr.ecr.us-east-1.amazonaws.com/support_org_provision:latest

# Organization Handler
buildorg:
	GOOS=linux go build -o dist/console/main ./cmd/console/org/ && zip -j dist/console/org_handler.zip dist/console/main && rm dist/console/main

deployorg: buildorg upload
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/console/org_handler.zip --function-name dev-console-handler-orgservice

# Scangroup Handler
buildscangroup:
	GOOS=linux go build -o dist/console/main ./cmd/console/scangroup/ && zip -j dist/console/scangroup_handler.zip dist/console/main && rm dist/console/main

deployscangroup: buildscangroup upload
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/console/scangroup_handler.zip --function-name dev-console-handler-scangroupservice

# Address Handler
buildaddress:
	GOOS=linux go build -o dist/console/main ./cmd/console/address/ && zip -j dist/console/address_handler.zip dist/console/main && rm dist/console/main

deployaddress: buildaddress upload
	aws lambda update-function-code --s3-bucket linkai-infra --s3-key frontend/lambdas/console/address_handler.zip --function-name dev-console-handler-addressservice
