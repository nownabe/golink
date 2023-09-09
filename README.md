# Golink (go/link)

A URL shortener for creating concise, memorable short links, suitable for organization-scoped use, such as in companies or schools.

https://github.com/nownabe/golink/assets/1286807/9337f9f8-b3de-40bd-879c-f21b8d441604

## Alternatives

Considering other options? Here are some similar platforms:

- [GoLinks® | Knowledge Discovery & Link Management Platform](https://www.golinks.io/)
- [Trotto - Open-Source Go Links](https://www.trot.to/)
- [tailscale/golink: A private shortlink service for tailnets](https://github.com/tailscale/golink)

Golinks has the following advantages compared to these alternatives:

- **Self-hosted**: You have full control of Golink backend.
- **Fully-managed**: Golink can be built on fully-managed infrastructure.
- **Easy to deploy**: You can complete to deploy Golink in just three simple steps.
- **Cost-effective**: You can get started with Golink at no cost.
- **No DNS configuration**: Redirect through the Golink Chrome extension.
- **Chrome extension with [Manifest V3](https://developer.chrome.com/docs/extensions/mv3/intro/)**: Manifest V3 is more secure than V2 and V2 will be end-of-life.

## Golink Origin

Golink originated from Google's internal short links.
If you are curious about history of go/, dive into the stories on these websites:

- [Golink: A private shortlink service for tailnets | Hacker News](https://news.ycombinator.com/item?id=33978767)
- [The GoLinks® Blog - The History of Go Links](https://www.golinks.com/blog/go-links-history/)
- [The Go Links Origin Story: Q&A with Benjamin Staffin · Trotto go links](https://www.trot.to/blog/2020/07/09/go-links-origin-story)

## Usage

### For General Users

1. Install the [Golink Chrome Extension](https://chrome.google.com/webstore/detail/golink/clecngohjeflemkblbfdfbjkjnigbjok).
2. Right-click the extension icon and select **Options**.
3. Input your Golink URL and then click the **Save** button.

## Setup for Administrators

### Prerequisites

- [New Google Cloud project](https://cloud.google.com/docs/get-started)
- [gcloud](https://cloud.google.com/sdk/docs/install) CLI

Additionally, you need to execute the following command:

```shell
gcloud auth login
```

### Configure Your Project

Set your project ID:

```shell
gcloud config set project <YOUR-PROJECT-ID>
```

### Deploy Applications

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

### Configure Identity-Aware Proxy

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

### Add Users

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

### Retrieve Your Golink URL

Determine your Golink URL with:

```shell
echo "https://$(gcloud app describe --format "get(defaultHostname)")"
```

Then notify your team members to enter this URL in Golink Chrome Extension Options. Enjoy using golinks!

### Distribute Golink extension to your organization

You can enforce Golink Chrome extension to be installed in your organization members' browsers.

1. Open https://admin.google.com and navigate to Devices > Chrome > Apps & extensions > Users & browsers.
2. Click the yellow plus button at the bottom right and then click "add from Chrome Web Store".
3. Enter `clecngohjeflemkblbfdfbjkjnigbjok` in the "View app by ID" textbox and click the "Select" button.
4. Set "Permissions and URL access" to "Allow all permissions".

<!--
4. Configure your Golink URL as JSON like follorings:

```js
{
  "golinkInstanceUrl": "https://your-golink.an.r.appspot.com"
}


```
-->