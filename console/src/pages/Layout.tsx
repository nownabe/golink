import React, { Suspense } from "react";
import { Outlet, Link, defer, useLoaderData, Await } from "react-router-dom";
import { Box } from "@mui/material";
import {
  AppBar,
  Toolbar,
  Typography,
  Drawer,
  Divider,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
} from "@mui/material";
import {
  Add as AddIcon,
  Link as LinkIcon,
  List as ListIcon,
} from "@mui/icons-material";
import zIndex from "@mui/material/styles/zIndex";

import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";

import ThemeRegistry from "@/components/ThemeRegistry";
import LinkComponent from "@/components/LinkComponent";
import client from "@/client";

const drawerWidth = 240;

export async function layoutLoader() {
  const email = (async () => {
    const resp = await client.getMe({});
    return resp.email;
  })();
  return defer({ email });
}

export default function Layout() {
  const { email } = useLoaderData() as ReturnType<typeof layoutLoader>;

  return (
    <ThemeRegistry>
      <Box sx={{ display: "flex" }}>
        <AppBar position="fixed" sx={{ zIndex: zIndex.drawer + 1 }}>
          <Toolbar>
            <Link
              to="/"
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
              <Suspense fallback="">
                <Await resolve={email}>{(email: string) => email}</Await>
              </Suspense>
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
              <ListItemButton LinkComponent={LinkComponent} href="/-/new">
                <ListItemIcon>
                  <AddIcon />
                </ListItemIcon>
                <ListItemText primary="Create New Golink" />
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton LinkComponent={LinkComponent} href="/-/">
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
            <Outlet />
          </Box>
        </Box>
      </Box>
    </ThemeRegistry>
  );
}
