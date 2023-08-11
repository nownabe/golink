import {
  Add as AddIcon,
  Link as LinkIcon,
  List as ListIcon,
} from "@mui/icons-material";
import {
  AppBar,
  Box,
  Divider,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
} from "@mui/material";
import zIndex from "@mui/material/styles/zIndex";
import type { Metadata } from "next";
import { headers } from "next/headers";
import Link from "next/link";

import ThemeRegistry from "../components/ThemeRegistry/ThemeRegistry";

export const metadata: Metadata = {
  title: "Golink",
};

const drawerWidth = 240;
const appEngineUserEmailHeader = "X-Goog-Authenticated-User-Email";

export function getUserEmail() {
  const authenticatedUserEmail = headers().get(appEngineUserEmailHeader);

  let email;

  if (!authenticatedUserEmail) {
    email = "";
  } else {
    email = authenticatedUserEmail[0].split(":")[1];
  }

  return email;
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const email = getUserEmail();

  return (
    <html lang="en">
      <body>
        <ThemeRegistry>
          <Box sx={{ display: "flex" }}>
            <AppBar position="fixed" sx={{ zIndex: zIndex.drawer + 1 }}>
              <Toolbar>
                <Link
                  href="/c/"
                  style={{
                    color: "#fff",
                    textDecoration: "none",
                    flexGrow: 1,
                  }}
                >
                  <Typography
                    variant="h6"
                    noWrap
                    component="h1"
                    sx={{
                      display: "flex",
                      color: "#fff",
                      alignItems: "center",
                      "&:active": { textDecoration: "none" },
                      gap: 1,
                    }}
                  >
                    <LinkIcon />
                    Golink
                  </Typography>
                </Link>
                <Typography variant="body1" component="span">
                  {email}
                </Typography>
              </Toolbar>
            </AppBar>
            <Drawer
              variant="permanent"
              anchor="left"
              sx={{
                width: drawerWidth,
                flexShrink: 0,
                "& .MuiDrawer-paper": {
                  width: drawerWidth,
                  boxSizing: "border-box",
                },
              }}
            >
              <Toolbar />
              <Divider />
              <List>
                <ListItem disablePadding>
                  <ListItemButton LinkComponent={Link} href="/c/-/new">
                    <ListItemIcon>
                      <AddIcon />
                    </ListItemIcon>
                    <ListItemText primary="Create New Golink" />
                  </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                  <ListItemButton LinkComponent={Link} href="/c/-/">
                    <ListItemIcon>
                      <ListIcon />
                    </ListItemIcon>
                    <ListItemText primary="My Golinks" />
                  </ListItemButton>
                </ListItem>
              </List>
            </Drawer>
            <Box sx={{ flexGrow: 1 }}>
              <Toolbar />
              <Box component="main" sx={{ p: 2 }}>
                {children}
              </Box>
            </Box>
          </Box>
        </ThemeRegistry>
      </body>
    </html>
  );
}
