export async function loadLayoutMiddleware(route) {
    try {
        let layout = route.meta.layout;
        let layoutComponent = await import(`@/views/layouts/${layout}.vue`);
        route.meta.layoutComponent = layoutComponent.default;
    } catch (e) {
        console.error("Error occurred in processing of layouts: ", e);
        console.log("Mounted default layouts Main");
        let layoutComponent = await import(`@/views/layouts/Error.vue`);
        route.meta.layoutComponent = layoutComponent.default;
    }
}
