import { Golink } from "@/gen/golink/v1/golink_pb";
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

export const metadata: Metadata = {
  title: "My Golinks | Golink",
};

const golinks: Golink[] = [
  {
    name: "mylink1",
    url: "https://mylink1.example.com",
    owners: ["myself@example.com"],
  },
  {
    name: "otherlink1",
    url: "https://otherlink1.example.com",
    owners: ["other@example.com"],
  },
];

export default function MyGolinks() {
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
