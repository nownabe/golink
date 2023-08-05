import { Code, ConnectError } from "@bufbuild/connect";
import Grid from "@mui/material/Unstable_Grid2";
import { Metadata, ResolvingMetadata } from "next";
import { redirect } from "next/navigation";

import client from "@/client";
import Error from "@/components/Error";

import Golink from "./Golink";

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

export default async function GolinkPage({ params }: Props) {
  let golink;

  try {
    const resp = await client.getGolink({ name: params.name });
    golink = resp.golink;
  } catch (e) {
    const err = ConnectError.from(e);

    if (err.code === Code.NotFound) {
      redirect(`/c/-/new?name=${params.name}`);
    }

    console.error(err);
    return (
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Error error={err} />
        </Grid>
      </Grid>
    );
  }

  if (!golink) {
    redirect(`/c/-/new?name=${params.name}`);
  }

  return (
    <Golink
      name={golink.name}
      url={golink.url}
      ownersFromServer={golink.owners}
    />
  );
}
