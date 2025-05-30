import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";
import useFetch from "../../../hooks/useFetch";
import { getUniversities } from "../../../api/admin/universities";

const Universities = () => {
  const { data, loading, error, fetchData } = useFetch(getUniversities);

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
      validate={config.validate}
      onSubmit={config.onSubmit}
    />
  );
};

export default Universities;
