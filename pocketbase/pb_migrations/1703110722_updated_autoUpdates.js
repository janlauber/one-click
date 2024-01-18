/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // remove
  collection.schema.removeField("pffoke7f")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "jmbmwt1e",
    "name": "policy",
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

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "weuvu4uj",
    "name": "pattern",
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
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "pffoke7f",
    "name": "sortType",
    "type": "select",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "semver",
        "timestamp"
      ]
    }
  }))

  // remove
  collection.schema.removeField("jmbmwt1e")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "weuvu4uj",
    "name": "regexPattern",
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
})
