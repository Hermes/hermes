hermes
======

An open-source distributed unlimited redundant backup solution.


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
