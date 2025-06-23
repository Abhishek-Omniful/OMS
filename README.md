# Order Management System (OMS)

This microservice-based Order Management System handles order ingestion, validation, and processing using AWS S3, SQS, Kafka, MongoDB, and internal IMS services.

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

| Component     | Technology                 |
|--------------|----------------------------|
| Language      | Go (Golang)               |
| Database      | MongoDB                   |
| Messaging     | AWS SQS, Kafka            |
| Storage       | AWS S3                    |





