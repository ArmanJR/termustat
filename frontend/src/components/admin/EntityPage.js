import { useState } from "react";
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
  canAdd=true,
  canEdit=true,
  canDelete=true,
}) => {
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
        onAdd={canAdd ? () => openDialog("add") : null}
      />
      {loading ? (
        <p></p>
      ) : error ? (
        <p>{error}</p>
      ) : (
        <EntityTable
          data={data}
          columns={tableColumns}
          onEdit={canEdit ? (item) => openDialog("edit", item) : null}
          onDelete={canEdit ? (item) => openDialog("delete", item) : null}
        />
      )}
      {(canAdd || canEdit || canDelete) && (
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
      )}
    </div>
  );
};

export default EntityPage;
