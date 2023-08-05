import { Golink } from "@/gen/golink/v1/golink_pb";
import { Button, Chip, Divider, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Metadata, ResolvingMetadata } from "next";
import DeletableOwners from "./DeletableOwners";
import { redirect } from "next/navigation";

const golinks: { [key: string]: Golink } = {
  mylink1: {
    name: "mylink1",
    url: "https://mylink1.example.com",
    owners: ["myself@example.com", "who@example.com"],
  },
  otherlink1: {
    name: "otherlink1",
    url: "https://otherlink1.example.com",
    owners: ["other@example.com", "who@example.com"],
  },
};

type Props = {
  params: {
    name: string;
  };
};

export async function generateMetadata(
  { params }: Props,
  _parent?: ResolvingMetadata
): Promise<Metadata> {
  return {
    title: `go/${params.name} | Golink`,
  };
}

export default function GolinkPage({ params }: Props) {
  const golink = golinks[params.name];
  if (!golink) {
    redirect(`/c/-/new?name=${params.name}`);
  }

  return (
    <Grid container spacing={2}>
      <Grid xs={12}>
        <Typography variant="h5" component="h2">
          go/{golink.name}
        </Typography>
      </Grid>
      <Grid xs={12}>
        <TextField label="URL" fullWidth value={golink.url} />
      </Grid>
      <Grid xs={12}>
        <Button variant="contained">Update</Button>
      </Grid>
      <Grid xs={12}>
        <Divider />
      </Grid>
      <Grid xs={12}>
        <Typography variant="h6" component="h3">
          Owners
        </Typography>
      </Grid>
      <Grid xs={12}>
        <DeletableOwners owners={golink.owners} />
      </Grid>
      <Grid xs={12}>
        <Typography variant="h6" component="h3">
          Add Owner
        </Typography>
      </Grid>
      <Grid xs={12}>
        <TextField
          label="email"
          fullWidth
          placeholder="new-owner@example.com"
        />
      </Grid>
      <Grid xs={2}>
        <Button variant="contained">Add Owner</Button>
      </Grid>
    </Grid>
  );
}
