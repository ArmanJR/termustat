import { useState, useEffect } from "react";
import styles from "./Faculties.module.css";
import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";
import useFetch from "../../../hooks/useFetch";
import { getUniversities } from "../../../api/admin/universities";
import { getFaculties } from "../../../api/admin/faculties";

const Faculties = () => {
  const { data: universities } = useFetch(getUniversities);
  const [selectedUniversity, setSelectedUniversity] = useState("");
  const { data: faculties, fetchData } = useFetch(
    selectedUniversity ? () => getFaculties(selectedUniversity) : null
  );

  useEffect(() => {
    if (selectedUniversity) {
      fetchData();
    }
  }, [selectedUniversity]);

  return (
    <>
      <h1>دانشکده‌ها</h1>

      <div className={styles.container}>
        <select
          className={styles.select}
          value={selectedUniversity}
          onChange={(e) => setSelectedUniversity(e.target.value)}
        >
          <option value="" disabled>
            دانشگاه را انتخاب کنید
          </option>
          {universities?.map((uni) => (
            <option key={uni.id} value={uni.id}>
              {uni.name_fa}
            </option>
          ))}
        </select>
      </div>

      {selectedUniversity && faculties && (
        <EntityPage
          key={selectedUniversity}
          title={`دانشکده‌های ${ universities?.find((u) => u.id === selectedUniversity)?.name_fa || "" }`}
          entityName={config.entityName}
          data={faculties}
          fetchData={fetchData}
          tableColumns={config.tableColumns}
          dialogFields={config.dialogFields}
          validate={config.validate}
          onSubmit={async ({ id, data, mode }) => {
            const fullData = {
              ...data,
              university_id: selectedUniversity,
            };
            await config.onSubmit({ id, data: fullData, mode });
          }}
        />
      )}
    </>
  );
};

export default Faculties;
