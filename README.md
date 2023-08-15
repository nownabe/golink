# Golink (go/link)

**Golink: Custom URL Shortener & Redirect**

Empower your web browsing experience with "Golink", a dynamic URL shortener. Turn cumbersome URLs into crisp and memorable short links prefixed with "go/", ideal for streamlined sharing and organization.

Key Features:

1. **Effortless URL Shortening**: Transform any webpage's URL into a concise `go/{link_name}` format with just a single click.
2. **Instant Redirection**: Navigate swiftly with `go/{link_name}` URLs, which redirect seamlessly to the original web URLs.
3. **Extension Popup Utility**: Generate short links for your active tab directly through the extension's popup interface.
4. **Private Server Deployment**: For a tailored experience, deploy your server on Google Cloud. This ensures a unique identity for your shortened links.
5. **Exclusive Access Control**: Designate access to chosen users or groups, ensuring that only those authorized can make use of your Golink server.
6. **Collaborative Server Capability**: Share your Golink server—and consequently, your short links—with peers, colleagues, and team members. This fosters seamless collaboration, ensuring team members can easily access vital resources.
7. **Optimize Business Workflow**: Perfect for corporate settings, leverage standardized short links to streamline work processes and enhance operational efficiency.

Note:

For an optimal experience, users are required to deploy the associated server on Google Cloud. Detailed guidelines are below.

**Make your URLs not just shorter, but smarter with Golink.**

## Usage

### For General Users

1. Install the [Golink Chrome Extension](https://chrome.google.com/webstore/detail/golink/clecngohjeflemkblbfdfbjkjnigbjok).
2. Right-click the extension icon and select **Options**.
3. Input your Golink URL and then click the **Save** button.

### Setup for Administrators

#### Prerequisites

- New Google Cloud project
- [gcloud](https://cloud.google.com/sdk/docs/install)

Additionally, you need to execute the following command:

```shell
gcloud auth login
```

#### Configure Your Project

Set your project ID:

```shell
gcloud config set project <YOUR-PROJECT-ID>
```

#### Deploy Applications

Clone this repository:

```shell
git clone https://github.com/nownabe/golink
cd golink
```

Run the deploy script. Replace `<REGION>` with one of [App Engine regions](https://cloud.google.com/about/locations#region).

```shell
./deploy.sh <REGION>
```

For instance:

```shell
./deploy.sh us-central1
```

#### Configure Identity-Aware Proxy

Begin by accessing the [Google Cloud Console](https://console.cloud.google.com/apis/credentials/consent) to set up the OAuth consent screen.

1. Choose User Type.
   - Opt for a user type based on your needs.
     For exclusive access to members of your organization, select `Internal`.
     Note: choosing `External` doesn't mean open access.
     Users can't access your Golink unless you grant explicit permission.
2. Enter App information
   - App name: `Golink`
   - User support email: Your email or a Google Group
   - Developer contact information: Your email or alternate contact emails
   - Finish by clicking **SAVE AND CONTINUE**
3. You don't have to configure scopes.

Proceed to [Identity-Aware Proxy](https://console.cloud.google.com/security/iap).
Turn on IAP for the App Engine app.
If you encounter an error status before enabling, you can safely disregard it at this time.

#### Add Users

To make Golink accessible to all members of your organization, execute:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member domain:<YOUR-COMPANY-DOMAIN>
```

If you prefer to grant access on an individual basis:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member user:<EMAIL>
```

You have the option to specify Google Groups too:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member group:<EMAIL>
```

Examples:

```shell
gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member domain:your-company.example.com

gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member user:user1@your-company.example.com

gcloud iap web add-iam-policy-binding \
  --role roles/iap.httpsResourceAccessor \
  --member group:group1@your-company.example.com
```

#### Retrieve Your Golink URL

Determine your Golink URL with:

```shell
echo "https://$(gcloud app describe --format "get(defaultHostname)")"
```

Then, enter this URL in Golink Chrome Extension Options. Enjoy using golinks!

## Alternatives

Considering other options? Here are some similar platforms:

- [GoLinks® | Knowledge Discovery & Link Management Platform](https://www.golinks.io/)
- [Trotto - Open-Source Go Links](https://www.trot.to/)
- [tailscale/golink: A private shortlink service for tailnets](https://github.com/tailscale/golink)
