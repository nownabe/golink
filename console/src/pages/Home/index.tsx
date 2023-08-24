import { Divider, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Suspense } from "react";
import { Helmet } from "react-helmet";
import { Await, defer, useLoaderData } from "react-router-dom";

import client from "@/client";
import GolinksList, { Loading } from "@/components/GolinksList";
import { Golink } from "@/gen/golink/v1/golink_pb";

export async function homeLoader() {
  const popularGolinks = (async () => {
    const resp = await client.listPopularGolinks({
      days: 28,
      limit: 20,
    });
    return resp.golinks;
  })();
  return defer({ popularGolinks });
}

export default function Home() {
  const { popularGolinks } = useLoaderData() as ReturnType<typeof homeLoader>;

  return (
    <>
      <Helmet>
        <title>Dashboard | Golink</title>
      </Helmet>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Typography variant="h5" component="h2">
            Golink Dashboard
          </Typography>
        </Grid>
        <Grid xs={12}>
          <Divider />
        </Grid>
        <Grid xs={12}>
          <Typography variant="h6" component="h3">
            Popular golinks
          </Typography>
        </Grid>
        <Grid xs={12}>
          <Suspense fallback={<Loading />}>
            <Await resolve={popularGolinks}>
              {(golinks: Golink[]) => <GolinksList golinks={golinks} />}
            </Await>
          </Suspense>
        </Grid>
      </Grid>
    </>
  );
}
