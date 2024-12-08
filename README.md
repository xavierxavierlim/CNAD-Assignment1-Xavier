Take note: I only started committing yesterday due to issues with GitHub permissions that prevented me from committing to my own repository. I was only able to resolve the problem yesterday, which is why there are currently only a few commits.

In the video, under the "View Profile" section, the membership details under "Membership Benefits" should show "Hourly Rate" instead of "Hourly Rate Discount". The video is uploaded to brightspace as I am not able to commit it to my repository here. 

Key Features:
- User Management 
    - Registration, login and view and update profile details
    - Membership tiers (Basic, Premium, VIP) with different hourly rates
    - Modify and cancel reservation dates, view reservation history and view invoices

- Vehicle Management
    - Display available vehicles based on specified dates
    - Provide an estimated cost based on the number of days (hours) and membership tier

- Reservation Management
    - Allow users to make reservations to available vehicles 
    - Allow users to modify and cancel reservations
  
- Billing Mangement
    - Allow users to pay for their reservations
    - Calculate the final cost based on teh duration of the reservation, membership tier and any valid promotion codes
    - Generate detailed invoices after payment is successfully made

Link to architecture diagram: https://github.com/user-attachments/assets/b3311fc1-f10c-43f5-83dd-57a6561a4bd3

If you are not able to view the architecture diagram, please view the PNG image or the Adobe XD file (CNAD_Xavier_ASG1_ArchitectureDiagram) inside my repository. 

The design of this microservices-based system focuses on achieving modularity, scalability, and efficient communication through RESTful APIs. The core functionality is distributed across several independent microservices, each responsible for a specific domain, ensuring clear separation of concerns. To maintain data consistency and reduce redundancy, all server components share a centralised database, with each microservice interacting with its designated table and any other relevant tableswithin the database. This design choice allows for streamlined data management while enabling microservices to operate independently, facilitating scalability, ease of maintenance, and flexibility for future enhancements.

On the client side, users and computers interact with the system through a collection of web-based user interfaces (UIs). These include the User Service Web UI, Vehicle Service Web UI, Reservation Service Web UI, and Payment Service Web UI, each catering to distinct aspects of the system's functionality. For instance, the User Service Web UI facilitates user management, such as login, register and profile updates; the Vehicle Service Web UI supports vehicle-related operations, such as the displaying of available vehicles; the Reservation Service Web UI handles reservation activities such as making bookings within specified dates; and the Payment Service Web UI manages billing and invoices. These UIs send and receive requests to backend services via an API Gateway, which acts as a central point for managing traffic, ensuring security, and routing requests to the appropriate microservices.

The server side is composed of several microservices, each corresponding to a specific functionality and table in the centralised database. These include the User Server (linked to the Users table), the Vehicle Server (connected to the Vehicles table), the Reservation Server (associated with the Reservations table), the Billing Server (linked to the Billing table), and the Promotion Server (associated with the Promotion table). Each microservice communicates with its corresponding table and other relevant tables in the shared database to retrieve, update, or insert data.

The API Gateway facilitates communication between the client-side UIs and the server-side microservices. It routes incoming requests from users and computers to the relevant microservices, ensuring a seamless interaction between the client and the server. Additionally, it also manages the APIs between the UIs and the client-side components, allowing the web UIs (such as the User Service Web UI, Vehicle Service Web UI, Reservation Service Web UI, and Payment Service Web UI) to interact with client-side services like User Management, Vehicle Management, Reservation Management, and Billing Management. These client-side services act as intermediaries, processing the requests from the UIs before passing them through the API Gateway to the corresponding server-side microservices. This tiered communication structure enhances modularity and separation of concerns, ensuring that the client-side logic remains distinct from the server-side operations.

The centralised database ensures consistency across the system, with each table representing a distinct data domain. By sharing a single database, the architecture minimizes duplication of data while maintaining clear boundaries between microservices. This design balances modularity and centralized data management, allowing the system to scale and adapt to changing requirements efficiently. Overall, this architecture is robust, flexible, and designed for scalability, supporting independent updates and maintenance of microservices.

To set up and run the microservices, follow these steps. Firstly, clone the repository into your local machine from GitHub and ensure that Go and MySQL are installed. Secondly, install the required packages, which are "github.com/go-sql-driver/mysql" for connecting to the database, and "github.com/joho/godotenv" for loading the environment variables from the env file. Thirdly, update the go.mod file by changing the module name to match the name of the folder containing all the project files. Fourthly, modify the main.go file by updating the controller path to reflect your folder structure. The path should include your folder name, followed by the name of the nested folder, and "client" (e.g., "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/client"). Lastly, in the client folder, update the model path in all files to point to your folder name, followed by the nested folder and "server" (e.g., "CNAD-ASSIGNMENT1-XAVIER/CNAD_ASG1/server"). Once these modifications are complete, you can execute the main.go file to run the entire application.
