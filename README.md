# Mygram

1. Copy the `.env.example` file:

   If you don't have an `.env` file, start by copying the provided `.env.example` file.
   
   ```
   cp .env.example .env
   ```

2. Open the .env file:

    Use your preferred text editor to open the .env file.

    ```
    nano .env
    ```
3. Update the Database Configuration:

    Locate the database configuration section in the .env file. Update it to match your provided database information:
    
    ```
    DB_HOST=your_database_host
    DB_PORT=your_database_port
    DB_DATABASE=your_database_name
    DB_USERNAME=your_database_username
    DB_PASSWORD=your_database_password
    DB_DEBUG=false
    ```
    Replace the placeholder values (`your_database_name`, `your_database_user`, `your_database_password`) with your actual database information. Set `DB_DEBUG` to `true` if you want to enable query logging to the console.


4. Import Postman Collection:
    Use Postman to import the provided collection file [Mygram Postman](https://github.com/fadilahonespot/mygram/blob/master/Mygram.postman_collection.json) to your postman for testing superindo API.
