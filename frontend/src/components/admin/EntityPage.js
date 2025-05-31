import { useState } from "react";
import EntityHeader from "./EntityHeader";
import EntityTable from "./EntityTable";
import EntityDialog from "./EntityDialog";

const EntityPage = ({
  title,
  entityName,
  data,
  fetchData,
  tableColumns,
  dialogFields,
  validate,
  onSubmit,
  onOpenDialog,
  canAdd=true,
  canEdit=true,
  canDelete=true,
}) => {
  const [dialogState, setDialogState] = useState({
    open: false,
    mode: null,
    entity: null,
  });

  const [localDialogFields, setLocalDialogFields] = useState(dialogFields);

  const openDialog = async (mode, entity = null) => {
    if (onOpenDialog) {
      await onOpenDialog({
        mode,
        entity,
        setDialogFields: setLocalDialogFields,
      });
    } else {
      setLocalDialogFields(dialogFields);
    }
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
      <EntityTable
        data={data}
        columns={tableColumns}
        onEdit={canEdit ? (item) => openDialog("edit", item) : null}
        onDelete={canEdit ? (item) => openDialog("delete", item) : null}
      />
      {(canAdd || canEdit || canDelete) && (
        <EntityDialog
          open={dialogState.open}
          mode={dialogState.mode}
          entityName={entityName}
          entity={dialogState.entity}
          fields={localDialogFields}
          validate={validate}
          onSubmit={handleSubmit}
          onClose={closeDialog}
        />
      )}
    </div>
  );
};

export default EntityPage;
