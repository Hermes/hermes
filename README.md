Hermes
======

An open-source distributed unlimited redundant backup solution.

With current backup storage systems, it is all too common to have a high investment, while getting a tiny amount storage, that often come along with privacy issues. Hermes aims to offers a free and open-source network of computers all around the world to provide a distributed and most importantly, anonymous backup solution. The system uses LZMA compression to provide a large amount of network redundancy without having dedicated servers. Files pushed to the network are also encrypted with AES, and split up into blocks to be sent all around the world.

Follow the development: https://trello.com/board/project-hermes/5197b968bb47d41233005620

## Usage

#### Non-Authenticated
Generate new vault key:

    hermes generate <output_file>

Load a pregenerated key:

    hermes load <vault_file>

#### Authenticated
Update vault.dat manifest / sync with network:

    hermes update

Lock active vault

    hermes lock
    
Push file to network
    
    hermes push <file>
    
Pull file from network

    hermes pull <file>

