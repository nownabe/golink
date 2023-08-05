import { ConnectError } from "@bufbuild/connect";
import { useCallback, useRef, useState } from "react";

import client from "@/client";

export function useUrl(name: string) {
  const ref = useRef<HTMLInputElement>(null);

  const [openSuccess, setOpenSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [updating, setUpdating] = useState(false);

  const onUrlUpdate = useCallback(() => {
    (async () => {
      if (!ref.current) {
        return;
      }

      setUpdating(true);
      try {
        await client.updateGolink({
          name: name,
          url: ref.current.value,
        });
        setOpenSuccess(true);
      } catch (e) {
        const err = ConnectError.from(e);
        console.error(err);
        setError(err.message);
      } finally {
        setUpdating(false);
      }
    })();
  }, [name, ref, setOpenSuccess, setUpdating, setError]);

  const onSuccessClose = useCallback(
    () => setOpenSuccess(false),
    [setOpenSuccess]
  );
  const onErrorClose = useCallback(() => setError(null), [setError]);

  return {
    urlRef: ref,
    urlUpdating: updating,
    onUrlUpdate,
    openUrlUpdateSuccess: openSuccess,
    onUrlUpdateSuccessClose: onSuccessClose,
    urlUpdateError: error,
    onUrlUpdateErrorClose: onErrorClose,
  };
}

export function useOwners(ownersFromServer: string[]) {
  const [owners, setOwners] = useState<string[]>(ownersFromServer);

  const [openRemoveSuccess, setOpenRemoveSuccess] = useState(false);
  const [removeError, setRemoveError] = useState<string | null>(null);
  const removeOwner = useCallback(
    (owner: string) => async () => {
      try {
        setOwners((owners: string[]): string[] =>
          owners.filter((o) => o !== owner)
        );
        await client.removeOwner({ owner });
        setOpenRemoveSuccess(true);
      } catch (e) {
        const err = ConnectError.from(e);
        console.error(err);
        setOwners((owners: string[]): string[] => [...owners, owner]);
        setRemoveError(err.message);
      }
    },
    [setOwners, setOpenRemoveSuccess, setRemoveError]
  );
  const onRemoveSuccessClose = useCallback(
    () => setOpenRemoveSuccess(false),
    [setOpenRemoveSuccess]
  );
  const onRemoveErrorClose = useCallback(
    () => setRemoveError(null),
    [setRemoveError]
  );

  const addRef = useRef<HTMLInputElement>(null);
  const [adding, setAdding] = useState(false);
  const [openAddSuccess, setOpenAddSuccess] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);
  const addOwner = useCallback(() => {
    (async () => {
      if (!addRef.current) {
        return;
      }

      console.log(addRef);
      console.log(addRef.current);
      console.log(addRef.current.value);
      const owner = addRef.current.value;
      console.log(owner);
      if (!owner || !owner.includes("@")) {
        setAddError("Owner email must be a valid email address");
        return;
      }

      setAdding(true);

      try {
        await client.addOwner({ owner });
        setOwners((owners: string[]): string[] => [...owners, owner]);
        setOpenAddSuccess(true);
      } catch (e) {
        const err = ConnectError.from(e);
        console.error(err);
        setAddError(err.message);
      } finally {
        setAdding(false);
      }
    })();
  }, [addRef, setAdding, setOwners, setOpenAddSuccess, setAddError]);
  const onAddSuccessClose = useCallback(
    () => setOpenAddSuccess(false),
    [setOpenAddSuccess]
  );
  const onAddErrorClose = useCallback(() => setAddError(null), [setAddError]);

  return {
    owners,
    removeOwner,
    openRemoveSuccess,
    onRemoveSuccessClose,
    removeError,
    onRemoveErrorClose,
    addOwner,
    addRef,
    adding,
    openAddSuccess,
    onAddSuccessClose,
    addError,
    onAddErrorClose,
  };
}
