import axiosInstance from "../axiosInstance";

export const getFaculties = async (universityId) => {
  const response = await axiosInstance.get(
    `/admin/universities/${universityId}/faculties`
  );
  return response;
};

export const addFaculty = async (facultyData) => {
  const response = await axiosInstance.post(
    "/admin/faculties",
    facultyData
  );
  return response;
};

export const deleteFaculty = async (id) => {
  const response = await axiosInstance.delete(
    `/admin/faculties/${id}`
  );
  return response;
};

export const editFaculty = async (id, facultyData) => {
  const response = await axiosInstance.put(
    `/admin/faculties/${id}`,
    facultyData
  );
  return response;
};
