# GraphQL Student Dashboard

## What it is

A browser dashboard for a Reboot01 learner account, showing identity, XP progress, audit ratio, skills, module selection, and recent project activity.

## How it works

`index.html` loads the login handler, which sends Basic authorization to the configured Reboot01 sign-in endpoint and saves the returned token. `profile.html` loads `profile-new.js`, which checks authentication, fetches the user and available modules, and loads dashboard sections. `api.js` defines GraphQL queries, `graphql.js` sends them with Bearer authorization, feature modules transform the returned data, and `charts.js` builds SVG charts. Module selection is saved and used to filter relevant XP/activity queries.

## Code map

- `index.html`, `js/login.js`: sign-in page and token creation.
- `profile.html`, `js/profile-new.js`: dashboard and loading coordination.
- `js/auth.js`: token-presence check and logout.
- `js/graphql.js`: authenticated GraphQL transport.
- `js/api.js`: queries for user, audits, XP, skills, modules, and activity.
- `js/modules/`: per-section loading and UI logic.
- `js/charts.js`: XP timeline, skill radar, audit graph, and detail displays.

## Run locally

Run these commands from this project’s folder.

Requires a modern browser, internet access, and a valid account on the configured Reboot01 service. No npm build is needed. With Python 3:

```bash
python3 -m http.server 8000
```

Open `http://localhost:8000` and sign in using your own account.

## Current behavior and limitations

The source uses `https://learn.reboot01.com/api/auth/signin` and `https://learn.reboot01.com/api/graphql-engine/v1/graphql`; this pass did not verify external availability or API compatibility. There is no local backend or demo dataset. The token and selected module ID live in localStorage; logout removes them. A token-presence check alone does not validate the token; requests still depend on server authorization. Login/API failures may also result from browser cross-origin restrictions.

## Try it

Sign in, select a module, inspect XP and activity, switch modules, reload to check selection persistence, then log out.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
