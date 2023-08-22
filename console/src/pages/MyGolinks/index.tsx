import { Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Helmet } from "react-helmet";

import { Suspense } from "react";
import { Await, defer, useLoaderData } from "react-router-dom";
import client from "@/client";
import { Golink } from "@/gen/golink/v1/golink_pb";
import GolinksList, { Loading } from "@/components/GolinksList";

export async function myGolinksLoader() {
  const golinks = (async () => {
    const resp = await client.listGolinks({});
    return resp.golinks;
  })();
  return defer({ golinks });
}

// TODO: Make it type safe (react-router-dom is not type safe now)
export default function MyGolinks() {
  const { golinks } = useLoaderData() as ReturnType<typeof myGolinksLoader>;

  return (
    <>
      <Helmet>
        <title>My Golinks | Golink</title>
      </Helmet>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Typography variant="h5" component="h2">
            My golinks
          </Typography>
        </Grid>
        <Grid xs={12}>
          <Suspense fallback={<Loading />}>
            <Await resolve={golinks}>
              {(golinks: Golink[]) => <GolinksList golinks={golinks} />}
            </Await>
          </Suspense>
        </Grid>
      </Grid>
    </>
  );
}
