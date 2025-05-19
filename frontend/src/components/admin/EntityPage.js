import { useState } from "react";
import { useTheme, useMediaQuery } from "@mui/material";
import useFetch from "../../hooks/useFetch";
import EntityHeader from "./EntityHeader";
import EntityTable from "./EntityTable";
import EntityDialog from "./EntityDialog";

const EntityPage = ({
  title,
  entityName,
  fetchFunction,
  tableColumns,
  dialogFields,
  validate,
  onSubmit,
}) => {
  const theme = useTheme();
  const isTablet = useMediaQuery(theme.breakpoints.down("md"));

  const { data, loading, error, fetchData } = useFetch(fetchFunction);

  const [dialogState, setDialogState] = useState({
    open: false,
    mode: null,
    entity: null,
  });

  const openDialog = (mode, entity = null) => {
    setDialogState({ open: true, mode, entity });
  };

  const closeDialog = () => {
    setDialogState((prev) => ({ ...prev, open: false }));
    setTimeout(() => {
      setDialogState({ open: false, mode: null, entity: null });
    }, 300);
  };

  const handleSubmit = async (formPayload) => {
    await onSubmit(formPayload);
    await fetchData();
  };

  return (
    <div>
      <EntityHeader
        title={title}
        entityName={entityName}
        onAdd={() => openDialog("add")}
      />
      {loading ? (
        <p></p>
      ) : error ? (
        <p>{error}</p>
      ) : (
        <EntityTable
          data={data}
          columns={tableColumns}
          isTablet={isTablet}
          onEdit={(item) => openDialog("edit", item)}
          onDelete={(item) => openDialog("delete", item)}
        />
      )}
      <EntityDialog
        open={dialogState.open}
        mode={dialogState.mode}
        entityName={entityName}
        entity={dialogState.entity}
        fields={dialogFields}
        validate={validate}
        onSubmit={handleSubmit}
        onClose={closeDialog}
      />
    </div>
  );
};

export default EntityPage;
