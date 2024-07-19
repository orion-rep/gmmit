# Git Authentication Tokens

## BitBucket

1. Login into your bickbucket account

2. Go to **Personal Settings** (Click the cog icon at the upper right corner).

3. Go to **ACCESS MANAGEMENT > App Password** on the left side menu.

4. Click on **Create app password**.

5. Add a label for the new password, and set the following permissions:

| Permission | Type |
| ---------- | ---- |
| Repositories | [x] Read |
| Repositories | [x] Write |
| Pull requests | [x] Read |
| Pull requests | [x] Write |

6. Click **Create**, and make sure to save the generated password on a safe place.

7. To generate the personal access token simply encode the string "<username>:<app_password>" using base64.

    i.e.

    ```bash
    echo 'luke-skywalker:fdJS()FJd8s9FJ8d9WZ0OkFUQkJERmN0tR34LRlzbjRXN1U2YlhlQTdOOUwzMDU4NEY0MAo=' | base64
    ```

    You should get something like: `bHVrZS1za3l3YWxrZXI6ZmRKUygpRkpkOHM5Rko4ZDlXWjBPa0ZVUWtKRVJtTjB0UjM0TFJsemJqUlhOMVUyWWxobFFUZE9PVXd6TURVNE5FWTBNQW89Cg==`

Check the official [BitBucket Docs](https://support.atlassian.com/bitbucket-cloud/docs/create-an-app-password/) for more details.