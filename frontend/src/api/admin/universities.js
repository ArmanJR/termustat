import axiosInstance from "../axiosInstance";

export const getUniversities = async () => {
  const response = await axiosInstance.get(
    "/admin/universities"
  );
  return response;
};

export const addUniversity = async (universityData) => {
  const response = await axiosInstance.post(
    "/admin/universities",
    universityData
  );
  return response;
};

export const deleteUniversity = async (id) => {
  const response = await axiosInstance.delete(
    `/admin/universities/${id}`
  );
  return response;
};

export const editUniversity = async (id, universityData) => {
  const response = await axiosInstance.put(
    `/admin/universities/${id}`,
    universityData
  );
  return response;
};
