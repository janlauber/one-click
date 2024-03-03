/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("xzzzptkywwwyxtr")

  // remove
  collection.schema.removeField("3ocuog7v")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "arsaxw3p",
    "name": "name",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "p1zgkntw",
    "name": "ingress_annotation_keys",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "tiwyifqr",
    "name": "ingress_default_annotations",
    "type": "json",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSize": 2000000
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "i3qv7cps",
    "name": "ingress_classes",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "asxceoga",
    "name": "ingress_hosts",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "150tacu4",
    "name": "storage_classes",
    "type": "text",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "min": null,
      "max": null,
      "pattern": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("xzzzptkywwwyxtr")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "3ocuog7v",
    "name": "ingressAnnotations",
    "type": "select",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "asdf"
      ]
    }
  }))

  // remove
  collection.schema.removeField("arsaxw3p")

  // remove
  collection.schema.removeField("p1zgkntw")

  // remove
  collection.schema.removeField("tiwyifqr")

  // remove
  collection.schema.removeField("i3qv7cps")

  // remove
  collection.schema.removeField("asxceoga")

  // remove
  collection.schema.removeField("150tacu4")

  return dao.saveCollection(collection)
})
