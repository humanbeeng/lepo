export type Data = {
	success: boolean;
};

export const actions = {
	default: async ({ request }) => {
		const formData = await request.formData();
		console.log(formData.get('email'));
		const success: Data = { success: true };

		return success;
	}
};
