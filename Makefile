localstack-start:
	docker-compose up -d

upload-csv:
	aws --endpoint-url=http://localhost:4566 s3api put-object --bucket sls-go-local-files-s3-bucket --key uploads/example.csv --body ./data/example.csv

invoke-s3: 
	sls invoke local --function s3ImportHandler --stage local --path ./event-payloads/s3-event.json --docker --docker-arg="--network sls-go_sls"
