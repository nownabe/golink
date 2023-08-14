# Golink (go/link)

**Golink: Custom URL Shortener & Redirect**

Empower your web browsing experience with "Golink", a dynamic URL shortener. Turn lengthy URLs into crisp and memorable short links prefixed with "go/" for efficient sharing and organization.

Key Features:

1. **Simple URL Shortening**: Convert any web page's URL into a concise `go/{link_name}` format with a single click.
2. **Instant Redirection**: Registered `go/{link_name}` URLs swiftly redirect to the original URLs, streamlining navigation.
3. **Extension Popup Utility**: Instantly create short links for your current tab directly from the extension's popup interface.
4. **Private Server Deployment**: To use this extension, deploy the server on Google Cloud, ensuring a personal touch to your shortened links.
5. **Exclusive Access Control**: Grant access to specific users, ensuring that only those permitted can utilize your server.
6. **Shareable Server Capability**: Share your server, and by extension, your shortened links with peers, colleagues, or team members.
7. **Optimize Business Workflow**: Perfect for corporate settings, harness common short links to streamline work processes and boost operational efficiency.

Note:

For an optimal experience, users are required to deploy the associated server on Google Cloud. Detailed guidelines are below.

**Make your URLs not just shorter, but smarter with Golink.**

## Usage

### General Users

1. Install [Golink Chrome Extension]()
2. Right click on extension icon and open **Options**
3. Enter the your Golink URL and click **Save** button

### Setup for Administrators

#### Prerequisites

- New Google Cloud project
- [gcloud](https://cloud.google.com/sdk/docs/install)

You also need to run `gcloud auth login`.

#### Configure Your Preference

Set project ID

```shell
gcloud config set project <your-project-id>
```

Set [region](https://cloud.google.com/about/locations#region)

```shell
gcloud config set compute/region <your-preferred-region> --quiet
```

#### Deploy Applications

Clone this repository.

```shell
git clone https://github.com/nownabe/golink
cd golink
```

Run deploy script. `region` must be one of [App Engine regions](https://cloud.google.com/about/locations#region).

```shell
./deploy.sh <region>
```

For example:

```shell
./deploy.sh us-central1
```

#### Configure Identity-Aware Proxy

Open [Google Cloud Console](https://console.cloud.google.com/apis/credentials/consent) and configure OAuth consent screen.

1. Choose user type. If you want to allow only members of your organization, choose `Internal`. Even if you choose `External`, any users cannot access your Golink until you explicitly allow them.
2. **App information**
  - App name: `Golink`
  - User support email: Your email or Google Group
  - Developer contact information: Your email or other contacts
  - Click **SAVE AND CONTINUE**
3. You don't have to configure scopes.

Go to [Identity-Aware Proxy](https://console.cloud.google.com/security/iap) and Enable IAP for App engine app.

Run the following command.

```shell
gcloud iap settings set \
	iap-settings.yaml \
	--resource-type=app-engine \
	--project="$(gcloud config get project)"
```

#### Add Users

If you want to make Golink available for all employees, run this command:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member domain:<your-company.example.com>
```

Instead, if you want allow each user to use Golink:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member user:<user1@your-company.example.com>
```

You can also specify Google Group:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member group:<group1@your-company.example.com>
```
