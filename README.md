<p align="center">
  <img src="resources/OMS%20-%20IMS%20WorkFlow%20Diagram.png" alt="OMS-IMS Workflow Diagram" width="100%">
</p>

# Order Management System (OMS)

This microservice-based Order Management System handles order ingestion, validation, and processing using AWS S3, SQS, Kafka, MongoDB, and internal IMS services.

---

## 📁 Project Directory Structure

OMS/
├── configs/
│ └── config.yaml
├── init/
│ └── initialize.go
├── localstack/
├── mycontext/
│ └── mycontext.go
├── pkg/
│ ├── controllers/
│ ├── helper/
│ │ ├── common/
│ │ ├── csvProcessor/
│ │ └── getLocalCSV/
│ ├── integrations/
│ │ ├── db/
│ │ ├── httpClient/
│ │ ├── kafka/
│ │ ├── redis/
│ │ ├── s3/
│ │ ├── sqs/
│ │ └── webhooks/
│ ├── middleware/
│ ├── models/
│ └── router/
├── public/
│ └── invalid_orders_20250705_063148.csv
├── resources/
│ ├── OMS - IMS WorkFlow Diagram.png
│ └── orders.csv
├── docker-compose.yaml
├── go.mod
├── go.sum
├── main.go
└── README.md


---

## 🚀 Workflow Overview

1. **CSV Upload**
   - The user uploads a CSV file containing orders via an API.
   - The backend uploads this file to an S3 bucket and returns a pre-signed path or acknowledges the upload.

2. **File Submission**
   - The user submits the **S3 file path** to a dedicated endpoint.
   - This triggers a validation process.
   - If the path and structure are valid, the path is published to an **SQS queue**.

3. **SQS Consumer**
   - A worker process listens to the SQS queue.
   - Upon receiving the S3 path:
     - Downloads the CSV.
     - Parses each row.
     - Extracts `hub_id` and `sku_id`.

4. **Validation with IMS**
   - For each row:
     - Calls IMS (Inventory Management System) API to validate `hub_id` and `sku_id`.
     - If valid:
       - Creates a new order in **MongoDB** with status `"onHold"`.
       - Publishes the order to a **Kafka topic**.

5. **Kafka Consumer**
   - Listens for new order events.
   - For each order:
     - Calls IMS again to check **inventory availability**.
     - If inventory is sufficient:
       - IMS deducts the requested quantity.
       - Returns a success response.
     - On receiving a `true` status from IMS:
       - **Updates** the order status in MongoDB to `"new Order"`.

---

## 📦 Tech Stack

| Component     | Technology     |
|---------------|----------------|
| Language      | Go (Golang)    |
| Database      | MongoDB        |
| Messaging     | AWS SQS, Kafka |
| Storage       | AWS S3         |

---

## 🛠️ Run Commands

```sh
# Create S3 Bucket in LocalStack
docker run --rm -it --network="host" \
  -e AWS_ACCESS_KEY_ID=test \
  -e AWS_SECRET_ACCESS_KEY=test \
  -e AWS_DEFAULT_REGION=us-east-1 \
  amazon/aws-cli \
  --endpoint-url=http://localhost:4566 \
  s3 mb s3://oms-temp-public

# Upload CSV to Local S3 Bucket
docker run --rm -it --network="host" \
  -e AWS_ACCESS_KEY_ID=test \
  -e AWS_SECRET_ACCESS_KEY=test \
  -e AWS_DEFAULT_REGION=us-east-1 \
  -v "${PWD}:/data" \
  amazon/aws-cli \
  s3 cp /data/orders.csv s3://oms-temp-public/orders.csv \
  --endpoint-url=http://localhost:4566

# Set environment variables and run OMS
$env:AWS_ACCESS_KEY_ID = "test"
$env:AWS_SECRET_ACCESS_KEY = "test"
$env:AWS_REGION = "us-east-1"
$env:AWS_S3_ENDPOINT = "http://localhost:4566"
$env:LOCAL_S3_BUCKET_URL = "localhost:4566"
$env:LOCALSTACK_ENDPOINT = "http://localhost:4566"
$env:LOCAL_SQS_ENDPOINT = "http://localhost:4566"
$env:CONFIG_SOURCE = "local"

go run main.go
