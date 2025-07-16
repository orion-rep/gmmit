# Gmmit Settings

This document explains the environment variables used to configure the application.
You can set these variables in your terminal for a single session or define them in a `.gmenv` file located in the root of your project for them to be loaded automatically.

## **Core Configuration**

This section covers the essential variables required for the application's primary functions, such as connecting to the generative AI service.

### **`GMMIT_API_KEY`**

This is a **mandatory** environment variable that holds your unique API key for accessing the generative AI service. The application cannot function without a valid API key.

* **Purpose:** To authenticate with the generative AI service.
* **Type:** String
* **Required:** Yes
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_API_KEY="your_unique_api_key_here"
        ```

  * **.gmenv File:**

        ```
        GMMIT_API_KEY="your_unique_api_key_here"
        ```

### **`GMMIT_MODEL`**

This optional environment variable allows you to specify which generative AI model you want to use. If you do not set this variable, the application will automatically use a default model.

* **Purpose:** To select the desired generative AI model.
* **Type:** String
* **Required:** No
* **Default Value:** `gemini-2.0-flash-lite`
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_MODEL="gemini-2.5-pro"
        ```

  * **.gmenv File:**

        ```
        GMMIT_MODEL="gemini-2.5-pro"
        ```

## **Git Integration Configuration**

These variables are used to authenticate with Git hosting providers. They are necessary only if you want the application to automatically create Pull Requests on your behalf.

### **`GMMIT_BB_USER`**

This variable holds your Bitbucket username. It is required for creating Pull Requests on Bitbucket.

* **Purpose:** To authenticate with the Bitbucket API.
* **Type:** String
* **Required:** No (Yes, if creating PRs on Bitbucket)
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_BB_USER="your_bitbucket_username"
        ```

  * **.gmenv File:**

        ```
        GMMIT_BB_USER="your_bitbucket_username"
        ```

### **`GMMIT_BB_PASS`**

This variable holds your Bitbucket App Password.

For security reasons, this should be an [App Password](https://support.atlassian.com/bitbucket-cloud/docs/create-and-use-app-passwords/) with "Pull Requests" write permissions, not your main account password.

* **Purpose:** To authenticate with the Bitbucket API.
* **Type:** String
* **Required:** No (Yes, if creating PRs on Bitbucket)
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_BB_PASS="your_bitbucket_app_password"
        ```

  * **.gmenv File:**

        ```
        GMMIT_BB_PASS="your_bitbucket_app_password"
        ```

### **`GMMIT_GH_USER`**

This variable holds your GitHub username. It is required for creating Pull Requests on GitHub.

* **Purpose:** To authenticate with the GitHub API.
* **Type:** String
* **Required:** No (Yes, if creating PRs on GitHub)
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_GH_USER="your_github_username"
        ```

  * **.gmenv File:**

        ```
        GMMIT_GH_USER="your_github_username"
        ```

### **`GMMIT_GH_PASS`**

This variable holds your GitHub Personal Access Token.
This should be a [Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) with the "repo" scope to allow for PR creation.

* **Purpose:** To authenticate with the GitHub API.
* **Type:** String
* **Required:** No (Yes, if creating PRs on GitHub)
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_GH_PASS="your_github_personal_access_token"
        ```

  * **.gmenv File:**

        ```
        GMMIT_GH_PASS="your_github_personal_access_token"
        ```

## **Feature Customization**

This section covers optional variables that allow you to customize the output of specific features, like the format of generated commit messages.

### **`GMMIT_COMMIT_PATTERN`**

This optional environment variable defines the structure for the generated Git commit messages, ensuring they adhere to a specific convention, such as the "Conventional Commits" standard.

* **Purpose:** To set a custom format for generated commit messages.
* **Type:** String
* **Required:** No
* **Default Value:** `"<type>[optional scope]: <description> (#<ticket-id>)"`
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_COMMIT_PATTERN="<type>: <description> TICKET-<ticket-id>"
        ```

  * **.gmenv File:**

        ```
        GMMIT_COMMIT_PATTERN="<type>: <description> TICKET-<ticket-id>"
        ```

## **Advanced Configuration**

These variables allow you to fine-tune the application's behavior for network resilience and debugging.

### **`GMMIT_MAX_RETRIES`**

This optional environment variable determines the maximum number of times the application will attempt to resend a request to the AI service.
If it fails with a specific type of error (a 500 Internal Server Error).

* **Purpose:** To define the maximum number of retries for failed AI requests.
* **Type:** Integer
* **Required:** No
* **Default Value:** `5`
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_MAX_RETRIES="10"
        ```

  * **.gmenv File:**

        ```
        GMMIT_MAX_RETRIES="10"
        ```

### **`GMMIT_RETRY_DELAY`**

This optional environment variable sets the waiting time, in seconds, between consecutive retry attempts.

* **Purpose:** To define the delay between retry attempts.
* **Type:** Integer
* **Required:** No
* **Default Value:** `5`
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_RETRY_DELAY="10"
        ```

  * **.gmenv File:**

        ```
        GMMIT_RETRY_DELAY="10"
        ```

### **`GMMIT_DEBUG`**

This variable is used to control the verbosity of the application's logging.
When enabled, it prints detailed debugging information to the console, which can be very useful for troubleshooting issues or understanding the application's internal workings.

* **Purpose:** To enable or disable debug-level logging.
* **Type:** Boolean (`"true"` or `"false"`)
* **Required:** No
* **Default Value:** `"false"`
* **Examples:**
  * **Shell / Terminal:**

        ```bash
        export GMMIT_DEBUG="true"
        ```

  * **.gmenv File:**

        ```
        GMMIT_DEBUG="true"
        ```
