import { InputAdornment, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { useCallback, useRef, useState } from "react";
import { Helmet } from "react-helmet";
import {
  LoaderFunctionArgs,
  useLoaderData,
  useNavigate,
} from "react-router-dom";

import ProgressButton from "@/components/ProgressButton";
import Snackbar from "@/components/Snackbar";
import { validateGolinkName, validateUrl } from "@/validator";
import client from "@/client";
import { ConnectError } from "@bufbuild/connect";

type LoaderData = {
  name: string | null;
};

export async function newGolinkLoader({
  request,
}: LoaderFunctionArgs): Promise<LoaderData> {
  const url = new URL(request.url);
  const name = url.searchParams.get("name");
  return { name };
}

export default function NewGolink() {
  const navigate = useNavigate();
  const { name } = useLoaderData() as LoaderData;

  const nameRef = useRef<HTMLInputElement>(null);
  const urlRef = useRef<HTMLInputElement>(null);
  const [creating, setCreating] = useState(false);
  const [openSuccess, setOpenSuccess] = useState(false);
  const [error, setError] = useState<string | null>("");

  const onClick = useCallback(() => {
    (async () => {
      if (!nameRef.current || !validateGolinkName(nameRef.current.value)) {
        setError("Invalid name");
        return;
      }
      if (!urlRef.current || !validateUrl(urlRef.current.value)) {
        setError("Invalid URL");
        return;
      }
      setCreating(true);

      try {
        const resp = await client.createGolink({
          name: nameRef.current.value,
          url: urlRef.current.value,
        });
        navigate(`/${resp.golink!.name}`);
      } catch (e) {
        const err = ConnectError.from(e);
        console.error(err);
        setError(err.message);
      } finally {
        setCreating(false);
      }
    })();
  }, [nameRef, urlRef, setCreating, setError, navigate]);
  return (
    <>
      <Helmet>
        <title>Create New Golink | Golink</title>
      </Helmet>
      <Grid container spacing={2}>
        <Grid xs={12}>
          <Typography variant="h5" component="h2">
            Create new golink
          </Typography>
        </Grid>
        <Grid xs={12}>
          <TextField
            label="Golink Name"
            inputRef={nameRef}
            fullWidth
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">go/</InputAdornment>
              ),
            }}
            placeholder="new-link-name"
            defaultValue={name}
          />
        </Grid>
        <Grid xs={12}>
          <TextField label="URL" inputRef={urlRef} fullWidth />
        </Grid>
        <Grid xs={12}>
          <ProgressButton loading={creating} onClick={onClick}>
            Create
          </ProgressButton>
        </Grid>
        <Snackbar
          open={!!error}
          severity="error"
          onClose={() => setError(null)}
        >
          {error || ""}
        </Snackbar>
        <Snackbar
          open={openSuccess}
          severity="success"
          onClose={() => setOpenSuccess(false)}
        >
          Successfully created.
        </Snackbar>
      </Grid>
    </>
  );
}
