"use client";

import { ConnectError } from "@bufbuild/connect";
import { InputAdornment, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useRef, useState } from "react";

import client from "@/client";
import ProgressButton from "@/components/ProgressButton";
import Snackbar from "@/components/Snackbar";

const validateName = (name: string): boolean => {
  if (name === "") {
    return false;
  }

  if (name.startsWith("-") || name.endsWith("-")) {
    return false;
  }

  if (name.match(/^_+$/) || name.startsWith("__") || name.endsWith("__")) {
    return false;
  }

  return true;
};

const validateUrl = (url: string): boolean => {
  try {
    new URL(url);
    return true;
  } catch (e) {
    return false;
  }
};

export default function Form() {
  const searchParams = useSearchParams();
  const name = searchParams.get("name");
  const router = useRouter();

  const nameRef = useRef<HTMLInputElement>(null);
  const urlRef = useRef<HTMLInputElement>(null);
  const [creating, setCreating] = useState(false);
  const [openSuccess, setOpenSuccess] = useState(false);
  const [error, setError] = useState<string | null>("");

  const onClick = useCallback(() => {
    (async () => {
      if (!nameRef.current || !validateName(nameRef.current.value)) {
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
        router.push(`/c/${resp.golink!.name}`);
      } catch (e) {
        const err = ConnectError.from(e);
        console.error(err);
        setError(err.message);
      } finally {
        setCreating(false);
      }
    })();
  }, [nameRef, urlRef, setCreating, setError]);

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
      <Snackbar open={!!error} severity="error" onClose={() => setError(null)}>
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
  );
}
