import {
  Await,
  LoaderFunctionArgs,
  defer,
  useLoaderData,
  useNavigate,
} from "react-router-dom";
import { Suspense } from "react";
import { Helmet } from "react-helmet";
import { Typography, Divider } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Code, ConnectError } from "@bufbuild/connect";
import { CircularProgress } from "@mui/material";

import client from "@/client";
import { Golink } from "@/gen/golink/v1/golink_pb";
import Owners from "./Owners";
import UpdateForm from "./UpdateForm";
import DeleteButton from "./DeleteButton";

export async function editGolinkLoader({ params }: LoaderFunctionArgs) {
  const name = params.name;
  const golink = (async () => {
    try {
      const resp = await client.getGolink({ name });
      return resp.golink;
    } catch (e) {
      const err = ConnectError.from(e);
      if (err.code === Code.NotFound) {
        return null;
      }
      throw err;
    }
  })();
  return defer({ name, golink });
}

// TODO: Make title a link
// TODO: Make it type safe (react-router-dom is not type safe now)
export default function EditGolink() {
  const navigate = useNavigate();
  const { name, golink } = useLoaderData() as ReturnType<
    typeof editGolinkLoader
  >;

  return (
    <>
      <Helmet>
        <title>go/{name} | Golink</title>
      </Helmet>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Typography variant="h5" component="h2">
            go/{name}
          </Typography>
        </Grid>

        <Suspense fallback={<CircularProgress />}>
          <Await resolve={golink}>
            {(golink: Golink | null) => {
              if (!golink) {
                navigate(`/-/new?name=${name}`);
                return null;
              }

              return (
                <>
                  <UpdateForm golink={golink} />
                  <Grid xs={12}>
                    <Divider />
                  </Grid>
                  <Owners golink={golink} />
                  <Grid xs={12}>
                    <Divider />
                  </Grid>
                  <DeleteButton golink={golink} />
                </>
              );
            }}
          </Await>
        </Suspense>
      </Grid>
    </>
  );
}
