# triple-s
## For PUT:
    http://localhost:8080/{BucketName}
    http://localhost:8080/{BucketName}/{ObjectKey}
## For GET:
    http://localhost:8080/
    http://localhost:8080/{BucketName}/{ObjectKey}
## For DELETE:
    http://localhost:8080/{BucketName}
    http://localhost:8080/{BucketName}/{ObjectKey}

### Example:

>##### Scenario 1: Bucket Creation
>- A client sends a `PUT` request to `/{BucketName}` with the name `my-bucket`.
>- The server checks for the validity and uniqueness of the bucket name, then creates an entry in the bucket metadata storage (e.g., `buckets.csv`).
>- The server responds with `200 OK` and the details of the new bucket or an appropriate error message if the creation fails.

>##### Scenario 2: Listing Buckets
>- A client sends a `GET` request to `/`.
>- The server reads the bucket metadata storage (e.g., `buckets.csv`) and returns an XML list of all bucket names and metadata.
>- The server responds with a `200 OK` status code.

>##### Scenario 3: Deleting a Bucket
>- A client sends a `DELETE` request to `/{BucketName}` for the bucket `my-bucket`.
>- The server checks if `my-bucket` exists and is empty.
>- If the conditions are met, the bucket is removed from the bucket metadata storage (e.g., `buckets.csv`).


### Example Scenarios

>- **Scenario 1: Object Upload**
   >  - A client sends a `PUT` request to `/photos/sunset.png` with the binary content of an image.
>  - The server checks if the `photos` bucket exists, validates the object key `sunset.png`, and saves the file to `data/photos/sunset.png`.
>  - The server updates `data/photos/objects.csv` with metadata for `sunset.png` and responds with `200 OK`.

>- **Scenario 2: Object Retrieval**
   >  - A client sends a `GET` request to `/photos/sunset.png`.
>  - The server checks if the `photos` bucket exists and if `sunset.png` exists within the bucket.
>  - The server reads the file from `data/photos/sunset.png` and returns the binary content with the `Content-Type` header set to `image/png`.

>- **Scenario 3: Object Deletion**
   >  - A client sends a `DELETE` request to `/photos/sunset.png`.
>  - The server checks if the `photos` bucket exists and if `sunset.png` exists within the bucket.
>  - The server deletes `data/photos/sunset.png` and removes the corresponding entry from `data/photos/objects.csv`.
>  - The server responds with `204 No Content`.

## Usage
Your program must be able to print usage information.

Outcomes:

- Program prints usage text.

```
$ ./triple-s --help  
Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory
```

