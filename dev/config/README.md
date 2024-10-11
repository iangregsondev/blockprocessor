# Project Configuration

This project uses a structured approach to manage configuration settings across multiple applications. Each application’s configuration is contained within its own subdirectory under this directory, allowing for modular and organized configuration management.

## Configuration Structure

**config.yml**

* Each subdirectory under this directory represents an individual application. The name of the subdirectory is the name of the application.
* Inside each subdirectory, there is a config.yml file that contains all the available configuration items for that specific application.
* This YAML file defines various settings and parameters required by the application, such as service URLs, feature toggles, and other static configuration values.
* **Important:** Environment variables can only override values that are defined in the config.yml file. If a configuration key is not present in the YAML file, the corresponding environment variable will not be used.

**env.local**

* There are two levels of env.local files:
  1. A global env.local file located in this directory, which serves as the default environment configuration for all applications.
  2. An individual env.local file located in the root of each application’s subdirectory. This file takes precedence over the global env.local, allowing for application-specific environment overrides.
     
* These files are useful for specifying sensitive configuration details, such as API keys, database passwords, or other secrets that should not be committed to version control.

### Recommended Usage

* It is recommended to use the config.yml files for all standard configuration settings. This keeps non-sensitive configurations in a structured format.
* The env.local files should be reserved for sensitive data, such as API keys, secrets, or credentials that should not be stored in version-controlled files.
* Define all configuration keys in config.yml: Even if the values are placeholders or empty, this ensures that environment variables can be used to override them.

### Naming Conventions for Environment Variables

Environment variables specified in the env.local files follow a standard naming pattern:

* Uppercase Naming: All environment variables are in uppercase.
* Application Prefix: Each environment variable is prefixed with the application name, which is written in uppercase. For example, an app named bitcoinProcessor would have a prefix of BITCOINPROCESSOR.
* Section and Key Separation: The sections from the YAML file are converted into environment variable format by separating each nested level with an underscore (_).
* Any camel-cased words from the YAML files are converted to uppercase with underscores separating each word. For instance, apiKey becomes API_KEY and rpcUrl becomes RPC_URL.

_For example:_

Given the following config.yml structure:

```yaml
connection:
  rpc-url: "https://some-url/"
  api-key: "your-sensitive-api-key-here"
```

The corresponding environment variable names would be:
* BITCOINPROCESSOR_CONNECTION_RPC_URL
* BITCOINPROCESSOR_CONNECTION_API_KEY

### Why Use env.local?

The env.local files allow for secure and flexible configuration management:

* **Sensitive Data Handling**: You can store sensitive information like API keys, database credentials, or secret tokens in env.local files. These details are kept out of source control and are protected from accidental exposure.
* **Environment Variable Overrides**: The env.local files use standard environment variable syntax, allowing configuration to be easily adjusted by exporting variables from the shell or setting them directly in the file.

### Git Ignore Settings

* Both the global and individual env.local files are excluded from version control through .gitignore. This ensures that sensitive information does not get committed to the repository, providing an extra layer of security.

### Example Usage

1. Set an API Key in env.local:

    ```bash
    BITCOINPROCESSOR_CONNECTION_API_KEY=your-sensitive-api-key-here
    ```

2. Export the Environment Variable:

    ```bash
    export BITCOINPROCESSOR_CONNECTION_API_KEY=your-sensitive-api-key-here
    ```

The env.local files provide an efficient way to manage and override environment variables, supporting a secure configuration practice for all applications.
