<template>
  <div>
    <div class="panel is-primary">
      <div class="panel-heading m-0">Heading</div>
      <div class="panel-block">
        <nav class="breadcrumb is-right" aria-label="breadcrumbs">
          <ul>
            <li><a href="#">Bulma</a></li>
            <li><a href="#">Documentation</a></li>
            <li><a href="#">Components</a></li>
            <li class="is-active">
              <a href="#" aria-current="page">Breadcrumb</a>
            </li>
          </ul>
        </nav>
      </div>
      <div class="panel-block">
        <p class="control has-icons-left">
          <input class="input is-primary" type="text" placeholder="Search" />
          <span class="icon is-small is-left">
            <i class="fas fa-search"></i>
          </span>
        </p>
      </div>
      <a class="panel-block pl-6">
        <span class="panel-icon">
          <i class="fas fa-file-alt" aria-hidden="true"></i>
        </span>
        marksheet
      </a>
    </div>
  </div>
</template>

<script>
export default {};
</script>
