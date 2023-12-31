import { onMounted, Ref, ref } from "vue";

export const useAsync = <T = any>(
    handler: (...args: any[]) => Promise<T>,
    asap: boolean = true
) => {
    const loading = ref(false);
    const data = ref<T | null>(null) as Ref<T>;
    const error = ref<any>(null);

    const action = async (...args: any[]): Promise<void> => {
        loading.value = true;
        error.value = null;
        try {
            const result = await handler(...args);
            data.value = result;
        } catch (err) {
            error.value = err;
            throw err;
        } finally {
            loading.value = false;
        }
    };

    onMounted(async () => {
        if (asap) {
            await action();
        }
    });

    return {
        data,
        loading,
        error,
        action,
    };
};
