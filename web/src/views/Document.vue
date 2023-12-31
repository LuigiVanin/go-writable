<script lang="ts" setup>
// Core
import { onMounted, provide, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

// Libraries
import { BubbleMenu, EditorContent } from "@tiptap/vue-3";

// Helpers & Services
import { DocumentSource, Html2PdfExporter } from "@/helpers/pdf/export.ts";

// Composables
import { useWritableEditor } from "@/composables/editor/useWritableEditor.ts";
import { useFetchDocument } from "@/composables/api/useFetchDocument.ts";
import { useUpdateDocument } from "@/composables/api/useUpdateDocument.ts";

// Components
import EditorActions from "@/components/editor/EditorActions.vue";
import TextInput from "@/components/core/TextInput.vue";
import StatusSelect from "@/components/editor/StatusSelect.vue";
import Button from "@/components/core/Button.vue";

// Types
import { DocumentStatus } from "@/types/document";

// Assets
import arrowLeftIconUrl from "@/assets/icons/arrow-left.svg?url";
import cloudCheckIconUrl from "@/assets/icons/cloud-check.svg?url";

const route = useRoute();
const router = useRouter();

const { editor } = useWritableEditor({
    onBlur: () => {
        saveDocument();
    },
});
const { fetchDocument, document, fetchError } = useFetchDocument();
const { updateDocument, updateLoading } = useUpdateDocument();

const title = ref("");
const status = ref<DocumentStatus>("draft");
const isExporting = ref(false);

provide("editor", {
    editor,
});

onMounted(async () => {
    if (!route.params?.id) {
        router.push("/");
    }

    await fetchDocument(route.params?.id as string);

    if (fetchError.value || !document.value) {
        router.push("/");
        return;
    }

    editor.value?.commands.setContent(
        document.value?.text.replaceAll("\n", "<p></p>")
    );
    status.value = document.value?.status || "draft";
    title.value = document.value?.title;
});

const saveDocument = async () => {
    if (!document.value) {
        return;
    }

    await updateDocument(document.value.id, {
        status: status.value || "draft",
        title: title.value,
        text: editor.value?.getHTML() || "",
    });
};

const exportAsPsf = () => {
    isExporting.value = true;
    const pdfExporter = new Html2PdfExporter();
    if (!editor.value) return;

    const source: DocumentSource = {
        title: title.value,
        content: editor.value?.getHTML() || "",
    };

    pdfExporter.export(source, () => {
        isExporting.value = false;
    });
};

const handleExport = () => {
    if (isExporting.value) {
        return;
    }

    exportAsPsf();
};
</script>

<template>
    <div class="editor-container">
        <header>
            <Button
                :icon="arrowLeftIconUrl"
                size="lg"
                btn-type="simple"
                @click="router.back()"
            >
                Voltar
            </Button>

            <div class="editor-container__header">
                <Button :loading="isExporting" size="lg" @click="handleExport">
                    Exportar para PDF
                </Button>
                <Button
                    :icon="cloudCheckIconUrl"
                    :loading="updateLoading"
                    loading-text="Salvando..."
                    size="lg"
                    btn-type="no-border"
                    @click="saveDocument"
                >
                    Salvamento automático
                </Button>
            </div>
        </header>
        <TextInput v-model="title" @blur="saveDocument" />
        <div class="editor-container__actions">
            <EditorActions />
            <StatusSelect v-model="status" @select="saveDocument" />
        </div>
        <div class="editor-content__wrapper">
            <EditorContent :editor="editor" />
            <BubbleMenu
                v-if="editor"
                :editor="editor"
                :tippy-options="{
                    animation: 'fade',
                }"
            >
                <div class="bubble-menu">
                    <EditorActions
                        :actions="[
                            'Bold',
                            'Italic',
                            'Underline',
                            'Strikethrough',
                            'Heading 1',
                            'Heading 2',
                            'Heading 3',
                        ]"
                    />
                </div>
            </BubbleMenu>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.editor-container {
    width: 100%;
    @include flex(column, center, center);
    gap: $spacing_13;

    header {
        width: 100%;
        @include flex(row, center, space-between);

        .editor-container__header {
            @include flex(row, center, center);
            gap: $spacing_5;

            :deep(.button-main.button-main--no-border) {
                svg path {
                    stroke: $green_3;
                }
            }

            @media (max-width: $lg_mobile) {
                & {
                    display: none;
                }
            }
        }
    }

    .editor-container__actions {
        width: 100%;
        @include flex(row, start, space-between);
        @media (max-width: $tablet) {
            flex-wrap: wrap;
            gap: $spacing_6;

            :deep(> .input-main) {
                margin-left: auto;
            }
        }
    }

    .editor-content__wrapper {
        width: 100%;
        height: calc(
            100vh - $header_height - 38px - 40px - 40px - (3 * $spacing_13) -
                (2 * $spacing_20)
        );
        @media (max-width: $tablet) {
            height: calc(
                100vh - $header_height - 38px - 40px - 40px - (3 * $spacing_13) -
                    (2 * $spacing_20) - 40px
            );
        }
        overflow-y: auto;
        overflow-x: visible;
        border-inline: $spacing_0 solid $neutral_2;
        @include set-scrollbar;
    }

    :deep(.input-main) {
        height: 40px;

        input {
            font-weight: 600;
            font-size: $font_8;
        }
    }
}
</style>

<style lang="scss">
.tiptap p.is-empty::before {
    color: #adb5bd;
    content: attr(data-placeholder);
    float: left;
    height: 0;
    pointer-events: none;
}

.editor-content__wrapper > div {
    height: 100% !important;
}

.tiptap.ProseMirror {
    padding-block: $spacing_8;
    min-height: 500px;
    height: 100%;
    outline: none;

    strong {
        font-weight: 700;
    }

    p,
    h1,
    h2,
    h3 {
        line-height: 1.5;
        margin-bottom: $spacing_6;
        color: $neutral_11;
    }

    h1 {
        font-size: $font_10;
    }

    h2 {
        font-size: $font_9;
    }

    h3 {
        font-size: $font_8;
    }

    em {
        font-style: italic;
    }

    blockquote {
        border-left: 4px solid $neutral_7;
        background: $neutral_3;
        padding-left: $spacing_8;
        padding-block: $spacing_2;
    }

    ul {
        list-style: disc;
        padding-inline-start: $spacing_16;

        li::marker {
            color: $neutral_11 !important; /* Change the color */
        }
    }

    ol {
        list-style: decimal;
        padding-inline-start: $spacing_16;

        li::marker {
            color: $neutral_11 !important; /* Change the color */
        }
    }
}

.editor-content__wrapper #tippy-2,
.editor-content__wrapper #tippy-1 {
    height: 40px !important;
    & > div {
        // background: red;
    }

    .bubble-menu {
        padding: $spacing_3 $spacing_4;
        border: $spacing_0 solid $neutral_3;
        background: $neutral_2;
        border-radius: $border_r_md;
        // width: 550px !important;

        // padding: $spacing_6 $spacing_6;
        width: 100%;
        box-shadow: 0px 0px 6px -2px $neutral_7;
    }
}
</style>
