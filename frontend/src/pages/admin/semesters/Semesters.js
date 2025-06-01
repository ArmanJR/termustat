import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";
import useFetch from "../../../hooks/useFetch";
import { getSemesters } from "../../../api/admin/semesters";

const Semesters = () => {
  const { data, loading, error, fetchData } = useFetch(getSemesters);

  if (loading) return <></>;
  if (error) return <div>{error}</div>;

  return (
    <EntityPage
      title={config.title}
      entityName={config.entityName}
      data={data}
      fetchData={fetchData}
      tableColumns={config.tableColumns}
      dialogFields={config.dialogFields}
      onSubmit={config.onSubmit}
    />
  );
};

export default Semesters;
