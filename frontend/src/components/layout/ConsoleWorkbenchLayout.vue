<template>
  <section class="console-workbench">
    <aside v-if="$slots.rail" class="console-workbench-rail">
      <slot name="rail" />
    </aside>

    <section class="console-workbench-main">
      <div v-if="$slots.summary" class="console-workbench-summary">
        <slot name="summary" />
      </div>

      <div class="console-workbench-ledger">
        <slot name="main" />
      </div>

      <div v-if="$slots.footer" class="console-workbench-footer">
        <slot name="footer" />
      </div>
    </section>
  </section>
</template>

<style scoped>
.console-workbench {
  display: grid;
  grid-template-columns: minmax(15rem, 0.26fr) minmax(0, 1fr);
  gap: 0.9rem;
  max-width: 1360px;
  margin: 0 auto;
  min-height: calc(100vh - 64px - 4rem - 11.5rem);
  align-items: stretch;
}

.console-workbench-rail {
  position: sticky;
  top: 5.5rem;
  align-self: start;
  display: flex;
  max-height: calc(100vh - 7rem);
  flex-direction: column;
  gap: 0.9rem;
  overflow: auto;
  border: 1px solid rgba(198, 184, 157, 0.46);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.032), transparent 28%),
    linear-gradient(180deg, rgba(250, 247, 239, 0.76), rgba(237, 229, 212, 0.48));
  padding: 0.95rem;
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.28);
}

.console-workbench-main {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 0.9rem;
}

.console-workbench-summary {
  flex: 0 0 auto;
}

.console-workbench-ledger {
  position: relative;
  display: flex;
  min-height: 0;
  flex: 1 1 auto;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.52);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.console-workbench-ledger::before {
  content: '';
  position: absolute;
  inset: 0 auto 0 0;
  width: 2px;
  background: linear-gradient(180deg, transparent, #a73a2a 18%, #a73a2a 82%, transparent);
  opacity: 0.72;
  pointer-events: none;
}

.console-workbench-ledger :deep(.table-wrapper) {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  scrollbar-gutter: stable;
}

.console-workbench-ledger :deep(table) {
  width: 100%;
  min-width: max-content;
  display: table;
}

.console-workbench-footer {
  flex: 0 0 auto;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.44);
  padding: 0.85rem 1rem;
}

@media (max-width: 1279px) {
  .console-workbench {
    grid-template-columns: 1fr;
    min-height: 0;
  }

  .console-workbench-rail {
    position: relative;
    top: auto;
    max-height: none;
  }

  .console-workbench-ledger {
    min-height: 32rem;
  }
}

@media (max-width: 768px) {
  .console-workbench-ledger {
    min-height: 0;
    overflow: visible;
  }

  .console-workbench-ledger :deep(.table-wrapper) {
    overflow: visible;
  }
}
</style>
<style>
.dark .console-workbench-rail,
.dark .console-workbench-ledger,
.dark .console-workbench-footer {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .console-workbench-rail {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.06), transparent 28%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.86), rgba(17, 19, 15, 0.68));
}
</style>
