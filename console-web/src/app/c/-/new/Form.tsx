"use client";

import { Typography, TextField, InputAdornment, Button } from "@mui/material";
import { useSearchParams } from "next/navigation";
import Grid from "@mui/material/Unstable_Grid2";

export default function Form() {
  const searchParams = useSearchParams();
  const name = searchParams.get("name");

  return (
    <Grid container spacing={2}>
      <Grid xs={12}>
        <Typography variant="h5" component="h2">
          Create new golink
        </Typography>
      </Grid>
      <Grid xs={12}>
        <TextField
          label="Golink Name"
          fullWidth
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">go/</InputAdornment>
            ),
          }}
          placeholder="new-link-name"
          value={name}
        />
      </Grid>
      <Grid xs={12}>
        <TextField label="URL" fullWidth />
      </Grid>
      <Grid xs={12}>
        <Button variant="contained">Create</Button>
      </Grid>
    </Grid>
  );
}
