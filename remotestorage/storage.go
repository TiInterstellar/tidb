package remotestorage

// The whole picture of write data into s3.
//
// 1. Create database in Glue (maybe the same name with tidb?)
// 2. Fetch the table schema, and create the linking table in Glue.
// 3. Setup crawler to allow Glue fetch data from the s3.
// 4. Create catalog at Athena side to load from Glue.
// 5. Write data (in csv or parquet) into s3 and trigger the Glue crawler.
//
// After all these setups, we can read data from Athena now.
