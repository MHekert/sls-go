docker-start:
	docker-compose up -d

s3-upload-csv:
	aws --endpoint-url=http://localhost:4566 s3api put-object --bucket sls-go-local-123456789-files-s3-bucket --key uploads/example.csv --body ./data/example.csv

sls-invoke-s3:
	sls invoke local --function s3ImportHandler --stage local --path ./event-payloads/s3-event.json --docker --docker-arg="--network sls-go_sls"

dynamo-scan:
	aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name sls-go-local-data-table --region eu-central-1
