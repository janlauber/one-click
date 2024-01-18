/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "lnnspbrn",
    "name": "repository",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "w6ia8paw6verog6",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // remove
  collection.schema.removeField("lnnspbrn")

  return dao.saveCollection(collection)
})
