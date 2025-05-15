import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";

const Universities = () => {
  return (
    <EntityPage
      title={config.title}
      entityName={config.entityName}
      fetchFunction={config.fetchFunction}
      tableColumns={config.tableColumns}
      dialogFields={config.dialogFields}
      validate={config.validate}
      onSubmit={config.onSubmit}
    />
  );
};

export default Universities;
