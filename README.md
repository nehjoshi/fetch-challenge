# About
This Go project is a submission for the Backend Engineer take-home assessment for Fetch.

## Details
The ```types``` folder contains the structs for receipts and items. Note that the receipts struct contains an array of items. 

The ```handlers``` folder contains 3 files:
- ```home.go```: A test, "Hello World!" handler for a basic home route
- ```process_receipt.go```: Contains the handler for processing receipts. Includes several helper functions that help calculate points for each separate condition. This file also contains a Map that stores the points for each receipt in-memory (will not persist on application restart). This handler is responsible for handling requests to the ```/receipts/process``` route.
- ```get_points.go```: Contains the handler for fetching the points of a receipt using an ```id``` param. Note that this file accesses the shared variable ```Scores``` (defined in ```process_receipt.go```) for fetching points. This handler is responsible for handling requests to the ```/receipts/{id}/points``` route.

```main.go``` acts as the API entry point and defines all possible routes. By default, the application will run on port 5000.