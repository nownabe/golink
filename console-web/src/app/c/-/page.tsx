import { ConnectError } from "@bufbuild/connect";
import {
  Divider,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Typography,
} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import type { Metadata } from "next";
import Link from "next/link";

import client from "@/client";
import Error from "@/components/Error";

export const metadata: Metadata = {
  title: "My Golinks | Golink",
};

export default async function MyGolinks() {
  let golinks;
  try {
    const resp = await client.listGolinks({});
    golinks = resp.golinks;
  } catch (e) {
    const err = ConnectError.from(e);
    console.error(err);
    return (
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Error error={err} />
        </Grid>
      </Grid>
    );
  }

  return (
    <Grid container spacing={2}>
      <Grid xs={12}>
        <Typography variant="h5" component="h2">
          My golinks
        </Typography>
      </Grid>
      <Grid xs={12}>
        <List>
          {golinks.map((golink) => (
            <>
              <ListItem key={golink.name} disablePadding>
                <ListItemButton
                  LinkComponent={Link}
                  href={`/c/${golink.name}`}
                  sx={{ pl: 1 }}
                >
                  <ListItemText
                    primary={`go/${golink.name}`}
                    secondary={golink.url}
                  />
                </ListItemButton>
              </ListItem>
              <Divider />
            </>
          ))}
        </List>
      </Grid>
    </Grid>
  );
}
