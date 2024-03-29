docker-start:
	docker-compose up -d

s3-upload-csv:
	aws --endpoint-url=http://localhost:4566 s3api put-object --bucket sls-go-local-123456789-files-s3-bucket --key uploads/example.csv --body ./data/example.csv

sls-invoke-s3:
	sls invoke local --function s3ImportHandler --stage local --path ./event-payloads/s3-event.json --docker --docker-arg="--network host"

offline:
	sls offline start --noTimeout --stage local

dynamo-scan:
	aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name sls-go-local-data-table --region eu-central-1

dynamo-get:
	aws --endpoint-url=http://localhost:4566 dynamodb get-item --table-name sls-go-local-data-table --key '{"id":{"S":"$(id)"}}' --region eu-central-1

dynamo-count:
	aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name sls-go-local-data-table --region eu-central-1 --select "COUNT"

dynamo-delete-table:
	aws --endpoint-url=http://localhost:4566 dynamodb delete-table --table-name sls-go-local-data-table --region eu-central-1

dynamo-create-table:
	aws --endpoint-url=http://localhost:4566 dynamodb create-table --table-name sls-go-local-data-table --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --region eu-central-1

dynamo-clear:
	make dynamo-delete-table && make dynamo-create-table

mockery-install:
	go install github.com/vektra/mockery/v2@v2.21.1

mockery:
	~/go/bin/mockery --all

utest:
	go test ./...