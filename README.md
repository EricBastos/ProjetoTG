# LAUNCHING A STABLECOIN WITH CROSS-CHAININTEROPERABILITY

## This Repository that supports my Final Work of Undergraduation Study

### Steps to reproduce the results

- First, make sure you have Docker Engine, Docker CLI and Docker Compose installed. You can find the instructions to do so at https://docs.docker.com/compose/install/
- Clone this repository somewhere you want:
```
git clone https://github.com/EricBastos/ProjetoTG.git
```
- You might need to grant permission to execute the local.sh script. To do so, execute:
```
chmod +x local.sh
```
- Simply execute local.sh, and, if you have docker set up correctly, it'll create and start all the containers
necessary to test the code.
- P.S.: If you're using windows, you can reproduce the exact same steps above if you do everything inside WSL. If you're not
using WSL, you might have to manually execute the commands inside local.sh in your command line. It should work as long as 
docker is properly set up.

If everything worked correctly, you should have 6 containers running:

![Running containers](./docs/runningContainers.png "Running containers")

The API is available at localhost:8080, and a simple Frontend which you can use to visually interact with the basic
functionalities of the API is available at localhost:80. A RabbitMQ and a Postgres instance are deployed along with the system,
but they are not exposed to the host machine. You can access them by entering their respective containers, though, if you want
to check the database tables, or the current queued messages, for example.

There are some parameters that can be configured in the .env file, such as the RPC hosts for Sepolia and Mumbai network,
the StableCoin contract addresses for both networks, and the deployer wallet private and public keys for both networks as well.
Other parameters shouldn't really be changed unless you know what you are doing. If you change any parameters, be sure to run
the local.sh script again to compile and update everything.

### Deploying your own StableCoin smart contract

The .env file is filled with a disposable wallet which I filled with some testnet tokens just for the sake of testing this project.
You might want to use your own wallet and deploy your own smart contract for your own testing. To do so, follow these steps:

- Todo