"use client";

import { Divider, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";

import Owners from "./Owners";
import UpdateForm from "./UpdateForm";

type Props = {
  name: string;
  url: string;
  ownersFromServer: string[];
};

export default function Golink({ name, url, ownersFromServer }: Props) {
  return (
    <Grid container spacing={2}>
      <Grid xs={12}>
        <Typography variant="h5" component="h2">
          go/{name}
        </Typography>
      </Grid>
      <UpdateForm name={name} url={url} />
      <Grid xs={12}>
        <Divider />
      </Grid>
      <Owners ownersFromServer={ownersFromServer} />
    </Grid>
  );
}
