import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";

const Semesters = () => {
  return (
    <EntityPage
      title={config.title}
      entityName={config.entityName}
      fetchFunction={config.fetchFunction}
      tableColumns={config.tableColumns}
      dialogFields={config.dialogFields}
      onSubmit={config.onSubmit}
    />
  );
};

export default Semesters;
